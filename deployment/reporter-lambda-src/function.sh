function handler () {
  EVENT_DATA=$1
  echo "$EVENT_DATA" 1>&2;
  ./kosli environment report ecs $ENV -C $ECS_CLUSTER --owner $KOSLI_USER
}