# prerender-go

Prepare your SPA for searching system bots like a boss!

`prerender-go` creates html pages for your SPA and allow you to put htmls as files to 
Amazon S3 bucket or in file system. Use them as response when google|yandex|whatever 
crawler visits your site by simple rewrite or proxy_pass rule in nginx.

## Usage

### Configure

Create configuration for prerender and storage in `configuration/<stage>` folder. 
Default stage is `development`. How to configure you can guess by reading configs. ;) 
If not - ask us.

### Run

In docker launch

1. Build image

        make image

2. Run image

         docker run --rm \
          --name prerender \
          -e STAGE=development \
          -v "$(pwd)"/configuration/development:/app/configuration/development \
          prerender:latest

Running on host machine (google chrome must be installed)

    make go

## LICENCE

MIT License

Copyright (c) 2020 SpaceTab.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
