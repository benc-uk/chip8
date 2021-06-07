start:
  cls
  ld V2, #b
  ld F, V2
	drw V0, V1, 5

  ld V2, #c
  ld F, V2
  ld V0, #5
  ld V1, #5
	drw V0, V1, 5
  call d1
  
end:
  jp end

d1:
  ld V2, #1
  ld F, V2
  ld V0, #a
  ld V1, #a
  drw V0, V1, 5
  RET