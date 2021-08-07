# prepare

    mkdir custom-image-deploy && cd custom-image-deploy

# go project
    
    go mod init example.com/custom-image-deploy

# kubebuilder

init project

    kubebuilder init --domain example.com CustomImageDeploy

create api

    kubebuilder create api --group customimagedeploy --version v1 --kind CustomImageDeploy
