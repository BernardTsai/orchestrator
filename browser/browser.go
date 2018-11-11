package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"tsai.eu/orchestrator/controller/file"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

func renderPath(path string, version string) (result string, err error) {
	// list all instances
	instances, err := listInstances(path, version)
	if err != nil {
		return "", err
	}

	// render all instances
	result, err = renderInstances(instances)
	if err != nil {
		return "", err
	}

	return
}

//------------------------------------------------------------------------------

type occurence struct {
	Number int // number of occurences
	Index  int // index of entry
}

type occurences []occurence

func (o occurences) Len() int {
	return len(o)
}
func (o occurences) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
func (o occurences) Less(i, j int) bool {
	return o[i].Number < o[j].Number
}

//------------------------------------------------------------------------------

func renderInstances(instances []*file.InstanceInfo) (result string, err error) {
	results := []string{}

	// loop over all instances and determine result
	for _, info := range instances {
		dictionary := map[string]string{}

		for key, dep := range info.Dependencies {
			if dep.Type == "service" {
				depEndpoint, _ := file.DecodeEndpoint(dep.Endpoint)

				path := depEndpoint.Path
				val, err := renderPath(path, dep.Version)
				if err != nil {
					return "", err
				}
				dictionary[key] = val
			}
		}

		// render template
		template := info.Configuration.Template

		for key, value := range dictionary {
			template = strings.Replace(template, "{{"+key+"}}", value, -1)
		}
		results = append(results, template)
	}

	// determine result with highest number of occurences
	if len(results) == 0 {
		return "", errors.New("no data")
	}

	// calculate all md5s
	md5s := map[[16]byte]([]int){}
	for index, res := range results {
		md5sum := md5.Sum([]byte(res))

		_, found := md5s[md5sum]
		if !found {
			md5s[md5sum] = []int{}
		}
		md5s[md5sum] = append(md5s[md5sum], index)
	}

	// sort according to occurences
	ocs := occurences{}
	for _, value := range md5s {
		oc := occurence{
			Number: len(value),
			Index:  value[0],
		}

		ocs = append(ocs, oc)
	}
	sort.Sort(ocs)

	// check if we have a quorum
	if len(ocs) > 1 && ocs[0].Number == ocs[1].Number {
		return "", errors.New("no quorum")
	}

	// success
	result = results[ocs[0].Index]

	return
}

//------------------------------------------------------------------------------
func listInstances(path string, version string) (instances []*file.InstanceInfo, err error) {
	// determine path of component data directory
	directory := filepath.Join(file.ROOTDIR, path, file.DATADIR)

	// list all files in directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// filter files
	for _, f := range files {
		if f.Name() != file.COMPFILE {
			// load file
			instanceFilename := filepath.Join(directory, f.Name())
			instanceInfo, err := file.LoadInstanceInfo(instanceFilename)
			if err != nil {
				return nil, err
			}

			// filter version
			if instanceInfo.Version != version {
				continue
			}

			// filter inactive components
			if instanceInfo.State != model.ActiveState {
				continue
			}

			// append object to result
			instances = append(instances, instanceInfo)
		}
	}

	// success
	return
}

//------------------------------------------------------------------------------

func render(w http.ResponseWriter, r *http.Request) {
	var path string
	var version string
	var parts []string
	var result string

	path = r.URL.Path
	path = strings.TrimPrefix(path, "/render/")
	parts = strings.Split(path, ":")

	// check if the parameters are correct
	if len(parts) != 2 {
		w.Write([]byte("invalid request"))
		return
	}

	// process request
	path = parts[0]
	version = parts[1]

	result, err := renderPath(path, version)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// success
	w.Write([]byte(result))
}

//------------------------------------------------------------------------------

type info struct {
	Name      string
	Component *file.ComponentInfo
	Instances map[string]*file.InstanceInfo
	Children  []*info
}

//------------------------------------------------------------------------------

func readInfo(name string, path string) (i *info, err error) {
	instances := map[string]*file.InstanceInfo{}
	var component *file.ComponentInfo

	// list all instances
	dataPath := filepath.Join(path, file.DATADIR)
	files, err := ioutil.ReadDir(dataPath)
	if err == nil {
		// filter files
		for _, f := range files {
			if f.Name() != file.COMPFILE {
				// load file
				instanceFilename := filepath.Join(path, file.DATADIR, f.Name())
				instanceInfo, err := file.LoadInstanceInfo(instanceFilename)
				if err != nil {
					return nil, err
				}
				instances[f.Name()] = instanceInfo
			} else {
				compPath := filepath.Join(path, file.DATADIR, file.COMPFILE)
				component, err = file.LoadComponentInfo(compPath)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// list all children
	files, err = ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// filter files
	children := []*info{}
	for _, f := range files {
		if f.Name() != file.DATADIR {
			// load file
			childPath := filepath.Join(path, f.Name())
			childInfo, err := readInfo(f.Name(), childPath)
			if err != nil {
				return nil, err
			}
			children = append(children, childInfo)
		}
	}

	// success
	result := info{
		Name:      name,
		Component: component,
		Instances: instances,
		Children:  children,
	}

	return &result, nil
}

//------------------------------------------------------------------------------

func data(w http.ResponseWriter, r *http.Request) {
	// process request
	info, err := readInfo("", file.ROOTDIR)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// render object
	yaml, _ := util.ConvertToYAML(info)

	// success
	w.Write([]byte(yaml))
}

//------------------------------------------------------------------------------

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/render/", render)

	http.HandleFunc("/data/", data)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//------------------------------------------------------------------------------
