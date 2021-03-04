package main

import (
	"github.com/Just4Ease/axonrpc/internal/gengoaxonrpc"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetPrefix("protoc-gen-axonrpc: ")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error: reading input: %v", err)
	}

	var request plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &request); err != nil {
		log.Fatalf("error: parsing input proto: %v", err)
	}

	if len(request.GetFileToGenerate()) == 0 {
		log.Fatal("error: no files to generate")
	}

	plg, err := protogen.Options{}.New(&request)
	if err != nil {
		log.Fatalf("error: mounting axonrpc protogen plugin: %v", err)
	}

	for _, f := range plg.Files {
		if !f.Generate {
			continue
		}
		g := gengo.GenerateFile(plg, f)
		gengoaxonrpc.GenerateFileContent(plg, f, g)
	}
	plg.SupportedFeatures = gengo.SupportedFeatures

	if data, err = proto.Marshal(plg.Response()); err != nil {
		log.Fatalf("error: failed to marshal output proto: %v", err)
	}

	if _, err := os.Stdout.Write(data); err != nil {
		log.Fatalf("error: failed to write output proto: %v", err)
	}
}
