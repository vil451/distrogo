For example, you have two Node.js projects: one requires Node version 18.10, and the other needs version 20.17. 
Instead of installing both versions of Node on your host system in different directories, 
you can configure the environment to automatically select the correct Node version for each project. 
This way, you can run Node as usual, and the tool will automatically choose the right version based on the project. 
Below is a code example for such projects.

## Node 20
Go project dir ProjectNode20

`cd usecases/nodejs/ProjectNode20`

Run build project

`docker build -f Dockerfile -t myprojectnode20 .`

Run distrogo container

`distrogo create -n node20  -i myprojectnode20:latest`


Enter distrogo container

`distrogo enter node20`

Run the installation and other what you need

`npm install`

Run server to start 

`node --experimental-modules index.js &`

Run curl for view Node version

`curl http://localhost:3000`

```
Hello World! Running on Node.js version: 20.17.0
Resolved path for 'fs' module: node:fs
```

## Node 18

Go project dir ProjectNode18

`cd usecases/nodejs/ProjectNode18`

Run build project

`docker build -f Dockerfile -t myprojectnode18 .`

Run distrogo container

`distrogo create -n node18 -i myprojectnode18:latest`

Enter distrogo container

`distrogo enter node18`

Run the installation and other what you need

`npm install`

Run server to start

`node --experimental-modules index.js &`

Run curl for view Node version

`curl http://localhost:3000`

```
Hello World! Running on Node.js version: 18.10.0
import.meta.resolve is not available in this version of Node.js
```

