# Proyecto de Blockchain

Fecha: 10-10-2023

Versión: 1.0

----

## Descripción del proyecto

El proyecto consiste en la elaboración de sistema de cadena de bloques (Blockchain). El sistema está implementado mediante Golang y LevelDB.

## Características

### Block

El archivo define la estructura `Block` y varias funciones y métodos que operan en bloques. Estas funciones y métodos se utilizan para crear nuevos bloques, calcular el hash de las transacciones, serializar y deserializar bloques, y crear el bloque génesis (o bloque origen).

### Blockchain

 El archivo define la estructura Blockchain y varias funciones y métodos que operan en bloques. Estas funciones y métodos se utilizan para crear nuevos bloques, inicializar una cadena de bloques, agregar bloques a la cadena, encontrar transacciones no gastadas, encontrar salidas de transacciones no gastadas y encontrar salidas gastables (estas últimas permiten la correcta validación de los fondos disponibles para las transacciones).

### Proof

Estos métodos se utilizan para inicializar los datos de prueba de trabajo, ejecutar la prueba de trabajo y validar la prueba de trabajo.

### CLI

Este módulo define la estructura disponible del sistema para la interacción del usuario para la blockchain mediante la terminal.

### Wallet

El archivo establece dos llaves: Pública y Privada. La Privada es única y se usa como identificador, mientras que la Pública se comparte. La generación de direcciones se realiza a través de cálculos en una curva elíptica con ECDSA. La llave Privada se extrae de este proceso, y la llave Pública se deriva mediante algoritmos de hash. La dirección Pública final se obtiene combinando varios elementos y procesándolos mediante base 58. En general, el módulo determina y facilita la creación de direcciones con lógica Privada para funciones en el sistema.

## Funciones disponibles


### createblockchain

Para inicializar el sistema, se ejecuta el siguiente comando:

```$ go run main.go createblockchain -address "<DIRECCION>"```

Se reemplaza `<DIRECCION>` por la dirección que se desea definir (Se recomienda el nombre del usuario). La cadena se inicializa con una recompensa de 100 unidades para el usuario de dicha dirección.

### getbalance

```$ go run main.go getbalance -address "<DIRECCION>"```

Se reemplaza `<DIRECCION>` por la dirección del usuario sobre el cual se desea conocer la cantidad de unidades del usuario asociado a la dirección proveída.

### send

```$ go run main.go send -from "<DIRECCION_1>" -to "<DIRECCION_2>" -amount <CANTIDAD>```

Se reemplaza `<DIRECCION_1>` por la dirección desde la que se desea enviar unidades. Se reemplaza `<DIRECCION_2>` por la dirección hacia la que se desea recibir unidades. Se reemplaza `<CANTIDAD>` por el total de unidades que se desea transferir.

### printchain

```$ go run main.go printchain```

Imprime la cadena.

### createwallet

```$go run main.go createwallet```

Entrega una dirección de wallet, diferente cada vez que se ejecuta. Se almacena en la DB automáticamente.

### listaddresses

```$go run main.go listaddresses```

Entrega el listado de direcciones creadas en el sistema.

## Referencias y fuentes

Log Rocket: https://blog.logrocket.com/build-blockchain-with-go/

Golang company: https://www.golang.company/blog/build-a-blockchain-with-golang

Tensor Programming: https://www.youtube.com/watch?v=mYlHT9bB6OE&list=PLJbE2Yu2zumC5QE39TQHBLYJDB2gfFE5Q

FCC: https://www.freecodecamp.org/news/build-a-blockchain-in-golang-from-scratch/

Anthony GG: https://www.youtube.com/watch?v=oCm46sUILcs&list=PL0xRBLFXXsP6-hxQmCDcl_BHJMm0mhxx7

