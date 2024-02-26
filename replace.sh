#!/bin/bash

find pkg/proto  -type f | xargs sed -i '' 's/github\.com\/authzed\/spicedb\/pkg\/proto\/validate/github\.com\/envoyproxy\/protoc-gen-validate\/validate/g'
find pkg/proto  -type f | xargs sed -i '' 's/github\.com\/authzed\/spicedb\/pkg\/proto\/google\/api\/expr\/v1alpha1/google\.golang\.org\/genproto\/googleapis\/api\/expr\/v1alpha1/g'
find pkg/proto  -type f | xargs sed -i '' 's/v12\ \"github\.com\/authzed\/spicedb\/pkg\/proto\/authzed\/api\/v1\"/v12\ \"github\.com\/authzed\/authzed-go\/proto\/authzed\/api\/v1\"/g'


