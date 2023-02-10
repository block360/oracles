while [ -n "$1" ]
do
case "$1" in
--gofer)
cp -v $2 $(pwd)/result/share/gofer.conf
export GOFER_CONFIG=$2
export GOFER_PATH=$(pwd)/result-3/bin/gofer
shift
  ;;
--omnia) 
cp -v $2 $(pwd)/result/share/feed.conf
export OMNIA_CONFIG=$2
export OMNIA_VERBOSE=true
export OMNIA_PATH=$(pwd)/result-2/bin/omnia
shift
;;
--spire) 
cp -v $2 $(pwd)/result/share/spire.conf
export SPIRE_CONFIG=$2
export SPIRE_PATH=$(pwd)/result-3/bin/spire
export SPIRE_LOG_VERBOSITY=debug
shift
 ;;
--) shift
break ;;
*) echo "$1 is not an option";;
esac
shift
done

export SETZER_MIN_MEDIAN=1
$(pwd)/result/bin/install-omnia feed