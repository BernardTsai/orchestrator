Shell Commands
==============

orchestrator

Usage:
  orchestrator cmd options
          model reset
                load <filename>
                save <filename>
                show
          domain list
                 create <domain>
                 show <domain>
                 load <domain> <filename>
                 save <domain> <filename>
                 delete <domain>
          task list <task>
               create <domain> <task> ....
               load <domain> <filename>
               save <domain> <task> <filename>
               show <domain> <task>
               delete <domain> <task>
          event list <domain>
                create <domain> <event> <type> <source>
                load <domain> <filename>
                save <domain> <event> <filename>
                show <domain> <event>
                delete <domain> <event>
          template list <domain>
                   create <domain> <template> <type> ... (tbd)
                   load <domain> <filename>
                   save <domain> <template> <filename>
                   show <domain> <template>
                   delete <domain> <template>
          variant list <domain> <template>
                  create <domain> <template> <variant> <configuration>
                  load <domain> <template> <filename>
                  save <domain> <template> <variant> <filename>
                  show <domain> <template> <variant>
                  delete <domain> <template> <variant>
          dependency list <domain> <template>
                     create <domain> <template> <variant> <dependency> <type> <component> <version>
                     load <domain> <template> <dependency> <filename>
                     save <domain> <template> <variant> <dependency> <filename>
                     show <domain> <template> <variant> <dependency>
                     delete <domain> <template> <variant> <dependency>
          architecture list <domain>
                       create <domain> <architecture>
                       load <domain> <filename>
                       save <domain> <architecture> <filename>
                       show <domain> <architecture>
                       delete <domain> <architecture>
          service list <domain> <architecture>
                  create <domain> <architecture> <service>
                  load <domain> <architecture> <service>
                  save <domain> <architecture> <service> <filename>
                  show <domain> <architecture> <service>
                  delete <domain> <architecture> <service>
          setup list <domain> <architecture>
                create <domain> <architecture> <service> <setup> <version> <state> <size>
                load <domain> <architecture> <service> <filename>
                save <domain> <architecture> <service> <setup> <filename>
                show <domain> <architecture> <service> <setup>
                delete <domain> <architecture> <service> <setup>
          component list <domain>
                    show <domain> <component>
                    load <domain> <filename>
                    save <domain> <component> <filename>
                    delete <domain> <component>
                    scale-out <domain> <component> <version>
                    scale-in <domain> <component> <version>
          instance list <domain> <component>
                   show <domain> <component> <instance>
                   load <domain> <component> <filename>
                   save <domain> <component> <instance> <filename>
                   delete <domain> <component> <instance>
