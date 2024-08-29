# Compiler
GO := go

# Flags
LDFLAGS := -ldflags="-s -w"

# Directories
CMD_DIR := cmd
SUBDIRS := enikos news247 newsbeast protoThema theFaq zougla

# Output directory
OUTPUT_DIR := bin

# Default target
all: build

# Build for the current OS
build: $(SUBDIRS)

# Build for Linux
linux: export GOOS=linux
linux: export GOARCH=amd64
linux: build

# Clean all binaries
clean:
	rm -rf $(OUTPUT_DIR)

# Clean a specific binary
clean-%:
	rm -f $(OUTPUT_DIR)/$*_*

# Create output directory
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# Build rules for each subdirectory
enikos: $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/*.go

news247:
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/main.go

newsbeast: $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/*.go

theFaq: $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/*.go

protoThema: $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/main.go

zougla: $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/*.go

.PHONY: all build linux clean clean-% $(SUBDIRS)


#This Makefile:
 #
 #Places all binaries in the root/bin directory.
 #Provides a clean target to delete all binaries.
 #Provides a clean-% target to delete a specific binary (e.g., make clean-subdir1).
 #Handles directories with different file requirements:
 #
 #subdir1, subdir3, and subdir5 build all .go files in their directories.
 #subdir2 builds specific files (main.go, extra1.go, extra2.go).
 #subdir4 only builds main.go.
 #
 #
 #
 #To use this Makefile:
 #
 #Place it in your project root (same level as the cmd directory).
 #Run make to build for your current OS.
 #Run make linux to build for Linux.
 #Run make clean to remove all built binaries.
 #Run make clean-subdir1 (replace subdir1 with the actual subdirectory name) to remove a specific binary.
 #
 #This Makefile will work for directories with different file requirements.
 #You can adjust the build rules for each subdirectory as needed. If you add new subdirectories or change the file requirements,
  #you'll need to update the corresponding build rules in the Makefile.

  #subdir2: $(OUTPUT_DIR)
     #	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$@_$(shell go env GOOS) $(CMD_DIR)/$@/main.go $(CMD_DIR)/$@/extra1.go $(CMD_DIR)/$@/extra2.go
