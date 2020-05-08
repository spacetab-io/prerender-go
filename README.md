# roastmap-go

Prepare your SPA for searching system bots like a boss!

`roastmap-go` creates html pages for your SPA and allow you to put htmls as files to 
Amazon S3 bucket or in file system. Use them as response when google|yandex|whatever 
crawler visits your site by simple rewrite or proxy_pass rule in nginx.

## Usage

### Configure

Create configuration for roastmap and storage in `configuration/<stage>` folder. 
Default stage is `development`. How to configure 

### Run

In docker launch

1. Build image

        make image

2. Run image

         docker run --rm \
          --name roastmap \
          -e STAGE=development \
          -v "$(pwd)"/configuration/development:/app/configuration/development \
          roastmap:latest

Running on host machine (google chrome must be installed)

    make go

## LICENCE

