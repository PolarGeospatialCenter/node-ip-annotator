# Kubernetes NODE IP Annotator
This project gets the IP address of the interface defined in the INTERFACE env variable and annotates the node with that IP using the ANNOTATION env variable. Your pod must have access to get and update nodes as well as host network access. Designed to be run as an init container to a CNI solution.
