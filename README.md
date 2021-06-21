# gke-deployment-restarter

For restarting those apps that you don't have time to debug. 

## TODO:
- Add github action that pushes to dockerhub
- Change the env var 'refresh' to be customizable


## How this works
Specify what namespace, and label/value key pairs you're looking to restart.
This will add an env var called 'refresh' to your deployment that creates a new UUID and rolls out a new deployment


## Helm values:

| Name                   | Default | Description |
|------------------------|---------|-------------|
| cronjob                |   gke-restarter      |   name of the cronjob in k8s          |
| roleName               |  gke-restarter       |   name of the clusterrole          |
| namespace              |  default       |  name of the namespace to deploy this to           |
| clusterRoleBindingName |   gke-restarter-rolebinding      |    name of the role cluster role binding         |
| serviceAccountName     |    gke-restarter-sa     |  name of the service account this will create           |
| schedule               |    "0 0 * * *"     |     crontab schedule when this will fire        |
| c_env.NAMESPACE        |  default       |  namespace of the deployment to restart           |
| c_env.LABEL_NAME       |  busybox       |  label of the deployment to restart           |
| c_env.LABEL_VALUE      |  busybox       |  value of previous label           |