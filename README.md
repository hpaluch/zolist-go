# ZoList - Zomato menu list


GAE Application (written in Go)

Goal is to list favorite menus using Zomato REST API as data source.

> WARNING! It is work in progress - currently it just list restaurants
> using remote REST API

This application currently demonstrates:
* How to 
  use [urlfetch API](https://cloud.google.com/appengine/docs/standard/go/urlfetch/reference) in Go/GAE to
  call [Zomato REST API](https://developers.zomato.com/documentation)
  json service and parse json 
  using [json.Unmarshall()](https://golang.org/pkg/encoding/json/).
  Just see [zolist/zoapi/zoapi.go](https://github.com/hpaluch/zolist-go/blob/master/zolist/zoapi/zoapi.go)
* How to [include your own packages](https://golang.org/pkg/encoding/json/)
  in your GAE program (it is real mess :-)

## Setup

To **properly** checkout source you must obey following structure:
```bash
cd 
mkdir -p src/github.com/hpaluch/
cd src/github.com/hpaluch/
git clone https://github.com/hpaluch/zolist-go.git
```

> REMEMBER! You must have parent directory structure
> exactly set to `src/github.com/hpaluch/` otherwise
> all local go imports like:
> ```go
> import (
>  ...
>	"github.com/hpaluch/zolist-go/zolist"
>  ...
> )
> ```
> Would fail!!!
> Please see discussion
> at https://cloud.google.com/appengine/docs/flexible/go/using-go-libraries


To Get Zomato API key:
* go to page https://developers.zomato.com/api
* click on `Generate API Key` button
* click on `Registrace` (= Register in English)
* fill in 
  * `Jméno a příjmení` (= Name and Surname)
  * `Emailová adresa`  (= e-mail address) 
  * `Heslo` (= password)
* click on `Zaregistrovat`
* login to your mail account
* click on confirmation mail from `Zomato` user.
* go back to https://developers.zomato.com/api
* click on `Generate API Key` button
* fill in
  * Phone & Companay/Blog URL
* and click on `Generate API Key`
* add your key to your `~/.bashrc` as:
```bash
export ZOMATO_API_KEY=your_key
```
* and source it:
```bash
source ~/.basrc
```



Install required components:

* Tested OS: `Ubuntu 16.04.3 LTS`, `x86_64`

* Install python 2.7 (or later 2.x) using:

  ```bash
  sudo apt-get install python2.7
  ```

* Download current Google Cloud SDK (formerly GAE SDK) from:
  https://cloud.google.com/appengine/docs/standard/go/download
  in my case
  https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-171.0.0-linux-x86_64.tar.gz 

* Unpack your archive somewhere for example under `/opt`
  (you might need root permission):

```bash
sudo mkdir /opt/gae
sudo chown $USER /opt/gae
tar xzf google-cloud-sdk-171.0.0-linux-x86_64.tar.gz -C /opt/gae
```
* Add newly created `/opt/google-cloud-sdk/` to your `PATH`,
  for example add this to your `~/.bashrc`:

```bash
export PATH=/opt/gae/google-cloud-sdk/bin:$PATH
```

* and reload environment using:

```bash
source ~/.bashrc
```

* add Go GAE plugin to your Google Cloud SDK:

```bash
gcloud components install app-engine-go
```

Create new application in GAE Dashboard:

* Go to your GAE Dashboard using this link:
  https://console.cloud.google.com/projectselector/appengine/create?lang=go
* Click on `Create` button
* Fill in unique _Project name_ (in my case `hp-zolist`)
* click on `Create` button
* confirm `us-central` as region
* click on `Cancel Tutorial` if it bugs you.

## Developing app

* to run this app locally use:
```bash
./run_dev.sh
```
* and go to URL: http://localhost:8080/
* to view cute Admin interface (something like "Dashboard Lite")
  use: http://localhost:8000

## Deploying app

For the first time you must register your Google Account to deploy app:

* configure your project ID (in my case `hp-zolist`):
```bash
gcloud config set project YOUR_APP_ID
```

* configure your Google Account for GAE:
```bash
gcloud config set account YOUR_GOOGLE_ACCOUNT
```
* login with your GAE account:
```bash
gcloud auth login
```
* new browser window should appear:
  * login or confirm selected account
  * allow required permissions for `Google Cloud SDK`
* you should see page with title "You are now authenticated with the Google Cloud SDK!"

And finally:
* to deploy app run script:
```bash
./deploy.sh
```


# Resources

* Official Go on GAE getting started:
  https://cloud.google.com/appengine/docs/standard/go/quickstart
* Understand variable declaration in Go:
  https://golang.org/ref/spec#Short_variable_declarations
* How to pass environment in GAE/Go:
  https://cloud.google.com/appengine/docs/standard/go/config/appref
* How to call JSON API:
  https://blog.alexellis.io/golang-json-api-client/
* Hot to import your libraries:
  https://cloud.google.com/appengine/docs/flexible/go/using-go-libraries
* Offtopic: Dynamic eval of variables in bash:
  http://www.tldp.org/LDP/abs/html/ivr.html

