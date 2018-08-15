
URL="http://localhost"
COUNT=$1
LOOPS=5
COOLDOWN=1000

echo "Building new version"
# build it so we are testing most current version
go build
echo "Done building"
echo ""

./crb -url=$URL -count=$COUNT -loops=$LOOPS -cooldown=$COOLDOWN -verbose -display