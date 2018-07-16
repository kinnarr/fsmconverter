# fsmconverter
A toml to  verilog fsm converter

## Introduction
The goal of this project is to simplify the tedious work of creating a finite state machine (FSM) in Verilog code by using the [TOML](https://github.com/toml-lang/toml "Tom's Obvious, Minimal Language") configuration file format.

## Install process
As this project is written in [Go](https://golang.org/ "The Go Programming Language") the Go development kit needs to installed on your system. How to do this depends on your OS or distribution.
For Fedora there is a helpful [article](https://developer.fedoraproject.org/tech/languages/go/go-installation.html) on how to get started. It should be useful for users of other distributions as well.
Once you have the GOPATH variable set and GOPATH/bin in your PATH as described in the above article,   simply execute
```sh
cd $GOPATH
go get -u -v github.com/kinarr/fsmconverter
cd src/github.com/kinarr/fsmconverter
go install
```
Typing 'fsmconverter' should now show you the help screen of fsmconverter.

## Usage
The states and their respective inputs, outputs and transitions are specified in one or several .toml files. These files are placed in one folder with an arbitrary sub-folder structure. 
Fsm converter reads in this directory(specified with the '--fsm-config-dir' flag) and either validates('fsmconverter validate') or prints the found states in an orderly way('fsmconverter prettyprint').

##File structure
### List states
```toml
[state.EXAMPLE]
```

### Conditions for state transition
```toml
[[state.EXAMPLE1.next.Example2.and.condition]]
condition1=3        #decimal works
condition2=0b00011  #so does binary
condition3=0o3      #octal
condition4=0x3      #and hexadecimal notation
```

#### Else/Default Transition
Sometimes you want to define a default transition that is taken only when the conditions of the other possible transitions are not met.
```toml
[state.EXAMPLE1.else.Example3]
```

### Inputs
Every condition needs to be defined in one of the toml files used, e.g. 'inputs.toml'
Same as in the outputs, the binary length of these variables has to be defined
```toml
[inputs]
condition1=2		#2'b00
condition2=5		#5'b00000
condition3=3		#3'b000
condition4=4		#4'b0000
```

### Outputs
Outputs need to be defined similarly to to the inputs; for every output variable a default value needs to be specified.
```toml
[outputs]
output1=3		#3'b000
output2=4		#4'b0000

[defaults.outputs]
output1=0		#3'b000
output2=15		#4'b1111
```
