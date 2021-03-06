# Veritas

Veritas is a cli for getting at Diego's truth.

## Downloading on a BOSH VM

For a linux build on a bosh vm:

```bash
 pushd $HOME
 wget http://onsi-public.s3.amazonaws.com/veritas -O ./veritas
 chmod +x ./veritas

 echo "export PATH=$PATH:$PWD" > veritas.bash
 echo "export DROPSONDE_ORIGIN=veritas" >> veritas.bash
 echo "export DROPSONDE_DESTINATION=localhost:3457" >> veritas.bash
 ./veritas autodetect >> veritas.bash
 ./veritas completions >> veritas.bash

 source ./veritas.bash
 popd
```

Once this is done, you simply need to `source ~/veritas.bash` when you log in again.

## Downloading on an OS X Workstation

For an OS X build (mainly for chugging logs locally):

```bash
  mkdir -p $HOME/bin

  pushd $HOME/bin

  wget http://onsi-public.s3.amazonaws.com/veritas-osx
  mv veritas-osx veritas
  chmod +x ./veritas

  popd
```
