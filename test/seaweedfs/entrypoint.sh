#!/bin/sh

isArgPassed() {
  arg="$1"
  argWithEqualSign="$1="
  shift
  while [ $# -gt 0 ]; do
    passedArg="$1"
    shift
    case $passedArg in
    $arg)
      return 0
      ;;
    $argWithEqualSign*)
      return 0
      ;;
    esac
  done
  return 1
}

case "$1" in

  'master')
  	ARGS="-mdir=/data -volumePreallocate -volumeSizeLimitMB=1024"
  	shift
  	exec /usr/bin/weed master $ARGS $@
	;;

  'volume')
  	ARGS="-dir=/data -max=0"
  	if isArgPassed "-max" "$@"; then
  	  ARGS="-dir=/data"
  	fi
  	shift
  	exec /usr/bin/weed volume $ARGS $@
	;;

  'filer')
  	ARGS=""
  	shift
  	exec /usr/bin/weed filer $ARGS $@
	;;
  *)
  	exec /usr/bin/weed $@
	;;
esac
