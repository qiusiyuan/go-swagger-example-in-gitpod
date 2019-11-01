# go-swagger-example-in-gitpod
A example go server using go-swagger [simple server tutorial](https://goswagger.io/tutorial/todo-list.html) with gitpod one-click environment setup.

## Content and Goal
This project will setup a go http server using go-swagger's [simple server tutorial](https://goswagger.io/tutorial/todo-list.html) with full implementation
using `map` data structure as suggested in the tutorial.
You may notice that this project only contains
* `swagger.yml` which is core file to go-swagger
* restapi `configure` file, which contains the main logic of http server's controller
* `.gitpod.yml` and `.gitpod.dockerfile`, which is gitpod's configuration for single-click setup

Goals of this project are
* Getting familiar of how go-swagger works
* Be able to play with a go-swagger http server
* Having gitpod Single-click setup for go-swagger
* Guides on how to setup a go-swagger http server from only a swagger.yml file

## Setup
There are 2 ways to run this project either using [gitpod for go-swagger-example](https://gitpod.io/#https://github.com/qiusiyuan/go-swagger-example-in-gitpod) or clone this repo and setup yourself.
I really suggest to using gitpod directly, this is simple and stable. If you would like to know more, please refer [gitpod doc](https://www.gitpod.io/docs/)

### 1. gitpod
1. Click this link [gitpod for go-swagger-example](https://gitpod.io/#https://github.com/qiusiyuan/go-swagger-example-in-gitpod), and follow the instruction of login if you haven't used gitpod before.
2. Wait for process done, it will take 1-2 min for the first time access.
3. Once the workspace is opened, terminals opened by gitpod will automatically
* Swagger generate server code from `swagger.yml`
* `go` get all the dependencies
* `go` install the server
* run server's binary on port 8765
Please wait for all the process done, you will see a message below to indicate this process is done.
```
2019/11/01 18:24:29 Serving a todo list application at http://127.0.0.1:8765
```
3. Now you can play with the server, choose `Terminal` on the top of the window and choose `New Terminal`.
![test3](https://user-images.githubusercontent.com/17970730/68048427-ef447800-fcb6-11e9-83dc-af22016737e3.png)

With opened new terminal, you can now play with it using curl.
* list all the items
```
curl -i localhost:8765
```
* add one item into list
```
curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
* delete an item (delete the first item for example `localhost:8765/1`)
```
curl -i localhost:8765/1 -X DELETE -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
* update an item
```
curl -i localhost:8765/2 -X PUT -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' -d '{"description":"go shopping"}'
```
Refer to [simple server tutorial](https://goswagger.io/tutorial/todo-list.html) for more info.

### 2. Clone and Setup
1. Firstly make sure you have go setup in your PC.
2. Download go-swagger binary into your $GOROOT directory.
Example of how I did it
``` bash
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url') \
  && curl -o $GOROOT/bin/swagger -L'#' "$download_url" \
  && chmod +x $GOROOT/bin/swagger
```
Other ways in [go swagger install](https://goswagger.io/install.html)
3. Git clone this repo to your `$GOPATH/src/`
4. Go into this project root directory
5. Swagger generate server code
```bash
swagger generate server -f swagger.yml
```
6. Install dependencies
```bash
go get ./...
```
7. Install the application
```bash
go install ./cmd/a-todo-list-application-server
```
The binary will be installed as `$GOPATH/bin/a-todo-list-application-server`
8. Run the application on port (ex. 8765)
```bash
$GOPATH/bin/a-todo-list-application-server --port 8765
```
9. Now you can access the service through `localhost:8765`
