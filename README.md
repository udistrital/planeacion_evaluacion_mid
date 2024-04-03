# planeacion_mid
Api Mid para el sistema de planeación universidad Distrital
## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

- Golang
- BeeGo



## Variables de Entorno
```
  PLANES_MID_HTTP_PORT: [Puerto de ejecución API]
  PLANES_SERVICE: [Servicio Api planes:crud]
```

NOTA: Las variables se pueden ver en el fichero conf/app.conf y están identificadas con SGA_MID_HTTP_PORT...
Ejecución del Proyecto

## Ejecución del proyecto
```
#1. Obtener el repositorio con Go
go get github.com/udistrital/planeacion_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/planeacion_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
PLANES_MID_HTTP_PORT=8080 PLANES_MID_SOME_VARIABLE=some_value bee run
```


## Ejecución Pruebas

Pruebas unitarias
```
# En Proceso
```

## Licencia

This file is part of planeacion_mid.

planeacion_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

planeacion_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with planeacion_mid. If not, see https://www.gnu.org/licenses/.
