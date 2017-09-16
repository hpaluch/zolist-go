# ZoList - Zomato menu list


GAE Application (written in go)
to list favorite menus - uses Zomato API as data source.

> WARNING! It is work in progress - currently it just list headers...


## Setup

Install required components:

* Tested OS: `Ubuntu 16.04.3 LTS`, `x86_64`

* Install python 2.7 (or later 2.x) using:

```bash
sudo apt-get install python2.7
```

* Download current Google Cloud SDK (formerly GAE SDK) from: \\
  https://cloud.google.com/appengine/docs/standard/go/download \\
  in my case \\
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

* Go to your GAE Dashboard using this link: \\
  https://console.cloud.google.com/projectselector/appengine/create?lang=go
* Click on `Create` button
* Fill in unique _Project name_ (in my case `hp-zolist` 
* click on `Create` button
* confirm `us-central` as region
* click on `Cancel Tutorial` if it bugs you.

Back in your app:

* open `app.yaml` and replace `hp-zolist` with your APP ID:

```
application: PUT_YOUR_APP_ID_HERE
```

## Developing app

* to run this app locally use:
```bash
./dev_run.sh
```
* and go to URL: ...



# Resources

* Official Go on GAE getting started:
  https://cloud.google.com/appengine/docs/standard/go/download
