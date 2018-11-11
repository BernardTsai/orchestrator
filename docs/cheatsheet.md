Orchestrator V1.0.0
===================

General Commands
----------------
help           
  list all commands
usage          
  list all commands in detail
clear          
  clear screen
comment <text>
  display <text>
exit           
  stop program

Model Commands
--------------  
model reset
  reset the contents of the model
model load
  load model from the file "model.json"
model save
  save model to the file "model.json"
model show
  display the contents of the model

Domain Commands
---------------
domain list
  list domains of the model
domain create <domain>
  create a domain with the specified name
domain show <domain>
  show the domain with the specified name
domain delete <domain>
  delete the domain with the specified name

Blueprint Commands
------------------
blueprint list <domain>
  list blueprints of the specified name
blueprint load <domain> <filename>
  load blueprint from a file into the domain with the specified name
blueprint show <domain> <blueprint>
blueprint delete <domain> <blueprint>
blueprint instantiate <domain> <blueprint>
