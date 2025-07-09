Minimum configuration
To use the module, the following needs to be specified:

```
module "bootstrap" {
    source = "github.com/rickliujh/kickstart-gogrpc//terraform/gcp/modules/bootstrap-gcp-account"

    state_file_region  = "region-for-state-file-bucket"
    state_file_bucket_name = "name-for-the-state-file-bucket"
}
```
