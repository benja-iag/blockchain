# Proyecto de Blockchain

Fecha: 10-10-2023

Versión: 1.0

----

## Descripción del proyecto

El proyecto consiste en la elaboración de sistema de cadena de bloques (Blockchain). El sistema está implementado mediante Golang y LevelDB.


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
