{{/*
  Copyright 2018 Franz Schmidt

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  		http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/}}
module fsm (
  clock,
  reset,
  {{ range $inputName, $inputLenght := .Inputs -}}
  {{$inputName}},
  {{end -}}
  state
  );
  {{$countStates := len $.States}}
  {{$binaryStateSize := getBinarySize $countStates }}
  parameter SIZE = {{ $binaryStateSize }};

  input wire clock, reset;
  {{ range $inputName, $inputLenght := .Inputs -}}
  input wire {{if gt $inputLenght 1}}[{{minus $inputLenght 1}}:0] {{end}}{{$inputName}};
  {{end -}}

  output reg [SIZE-1:0] state;

  parameter {{ range $index, $stateName := enumerateKeys .States }}{{upper $stateName}} = {{convertBinary $index $binaryStateSize}}{{if ne $index (minus $countStates 1)}}, {{end}}{{end}};

  reg [SIZE-1:0] next_state;

  always @ (state{{ range $inputName, $inputLenght := .Inputs }} or {{$inputName}}{{end}})
  begin
    next_state = {{convertBinary 0 $binaryStateSize}};
    case(state)
    {{- range $stateName, $state := .States}}
      {{$countCondition := 0}}
      {{ upper $stateName}} : {{range $nextName, $next := $state.Successors -}}
                {{if or $next.And $next.Or}}{{$countCondition := inc $countCondition}}if {{if $next.And}}({{conditionToString $next.And "&&"}}){{end -}}
                {{- if $next.Or}}({{conditionToString $next.Or "||"}}){{end}}{{end}} begin
                  next_state = {{upper $nextName}};
                end {{if or $next.And $next.Or}}else {{end}}{{end -}}{{if ne $countCondition 0}}begin
                  {{- range $elseName, $elseValue := $state.DefaultSuccessor}}
                  next_state = {{upper $elseName}};
                  {{- else}}
                  next_state = {{upper $stateName}};
                  {{- end}}
                end
                {{- end}}
    {{end -}}
      default : next_state = {{upper .Defaults.State}};
    endcase
  end

  always @ (posedge clock or posedge reset)
  begin
    if (reset == 1'b1) begin
      state <= {{upper .Defaults.State}};
    end else begin
      state <= next_state;
    end
  end
endmodule
