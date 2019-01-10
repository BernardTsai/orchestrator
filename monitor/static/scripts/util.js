var loadData = async(url) => {
  response = await fetch(url);
  text     = await response.text();

  return text
}

function dump(obj) {
  console.log(jsyaml.dump(obj))
}
