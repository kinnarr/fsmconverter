`timescale 1 ns / 10 ps // module timescale_check5;
module test;
	{{ range $outputName, $outputLenght := .Outputs -}}
	wire {{if gt $outputLenght 1}}[{{minus $outputLenght 1}}:0] {{end}}{{$outputName}};
	{{end -}}
	{{ range $inputName, $inputLenght := .Inputs -}}
  reg {{if gt $inputLenght 1}}[{{minus $inputLenght 1}}:0] {{end}}{{$inputName}};
  {{end -}}
	{{$countStates := len $.States}}
  {{$binaryStateSize := getBinarySize $countStates }}
  parameter SIZE = {{ $binaryStateSize }};

	reg clock;
  reg reset;

	wire [SIZE-1:0] state;
  wire state_set;

	fsm F0 (
		{{ range $inputName, $inputLenght := .Inputs -}}
	  .{{$inputName}} ({{$inputName}}),
	  {{end -}}
		.state (state),
		.state_set (state_set),
		.clock (clock),
		.reset (reset)
	);

  cu C0 (
		{{ range $outputName, $outputLenght := .Outputs -}}
		.{{$outputName}} ({{$outputName}}),
		{{end -}}
		.state (state),
		.state_set (state_set)
  );

	initial
	begin
		clock = 0;
		reset = 0;
		{{ range $inputName, $inputLenght := .Inputs -}}
	  {{$inputName}} = 0;
	  {{end -}}
	end

	always
		#1	clock = !clock;

	initial
	begin
		$dumpfile ("fsm.vcd");
		$dumpvars;
	end

	initial
	begin
		$display("\t\ttime,\tclock,\treset,\tstate,\tstate_set");
		$monitor("%d,\t%b,\t%b,\t%d,\t\t%b",$time, clock, reset, state, state_set);
	end

	initial
		#350 $finish; // TODO: edit me
	initial
	begin
		/* insert test code */

	end
endmodule
