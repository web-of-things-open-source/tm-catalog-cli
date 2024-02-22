# tm-catalog-cli
A CLI Client for ThingModel Catalogs

This software is **experimental** and may not be fit for any purposes. 


go run main.go list -s temperature --filter.manufacturer Siemens
go run main.go search -n 1000 "+schema\:manufacturer.schema\:name:siemens +temperature"

bleve search ../catalog.bleve "+schema\:manufacturer.schema\:name:*" -l 10000

./tm-catalog-cli createSi ; du -s ../cata*; bleve count ../catalog.bleve


go install github.com/blevesearch/bleve/v2/cmd/bleve@latest
rm -rf ../catalog.bleve
bleve create ../catalog.bleve
bleve index ../catalog.bleve ../tm-catalog2
bleve count ../catalog.bleve
bleve query ../catalog.bleve "schema\:manufacturer.schema\:name:siemens" -l 1000

./tm-catalog-cli createSi -r siemens-dev ->> geht nicht
./tm-catalog-cli createSi -r localfs2 ->> geht
