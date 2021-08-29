.data
  _f: .word 0
  _g: .word 5
  _h: .word 4
  _i: .word 2
  _j: .word 2 
.text
  #Carrega os valores das variáveis nos registradores MEM -> REG
  lw $r1, _g
  lw $r2, _h
  lw $r3, _i
  lw $r4, _j
  #Realiza o calculo
  bne $r3,$r4,Else  #se i!=j => Else
  add $r0,$r1,$r2   #f=g+h (se i!=j pula, se i==j executa) 
  j Exit            #salto incondicional para Exit
  Else:
  sub $r0,$r1,$r2   #se i!=j, f=g-h  
  Exit:
  #Retorna o valor do resultado para a memória REG -> MEM
  sw $r0, _f
  #Termina o programa
  li $r0, 10
