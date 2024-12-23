# Web Scraper API - TDL

Una API con una interfaz gráfica que permite realizar scraping de computadoras y notebooks en diferentes sitios web según parámetros personalizados.

## **Integrantes**

- **Martin Alejo Polese** - 106808
- **Lucas Grati** - 102676
- **Edgardo Francisco Saez** - 104896

---

## **Requisitos previos**

Asegúrate de tener instalados los siguientes componentes antes de ejecutar el proyecto:

- [Go](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/) con npm

---

## **Instalación y ejecución**

### **Frontend**

1. Posicionarse en el directorio correspondiente al frontend.
2. Ejecutar el siguiente comando para instalar las dependencias necesarias:

   ```bash
   npm install
   ```

3. Para levantar el servidor en modo desarrollo, ejecutar:

   ```bash
   npm start
   ```

   Esto iniciará la aplicación en [http://localhost:3000](http://localhost:3000).

### **Backend**

1. Posicionarse en el directorio correspondiente al backend.
2. Ejecutar el siguiente comando para iniciar el servidor:

   ```bash
   go run .
   ```

   El servidor estará disponible en el puerto 8080.

---

## **Endpoints de la API**

Todas las solicitudes a la API deben realizarse utilizando el método **GET**. Recuerda anteponer la URL base: [http://localhost:8080](http://localhost:8080).

### **General**

Scrapea notebooks de Mercado Libre, Frávega y FullH4rd:

```bash
/api/general
```

### **Mercado Libre**

Scrapea notebooks específicamente de Mercado Libre:

```bash
/api/mercadolibre
```

### **Frávega**

Scrapea notebooks específicamente de Frávega:

```bash
/api/fravega
```

### **FullH4rd**

Scrapea notebooks específicamente de FullH4rd:

```bash
/api/fullh4rd
```

---

## **Parámetros de consulta (Query Parameters)**

Puedes personalizar el scraping utilizando los siguientes parámetros:

### **RAM**:
- `MinRam`: Memoria RAM mínima (en GB).
- `MaxRam`: Memoria RAM máxima (en GB).

### **Tamaño de pantalla (pulgadas)**:
- `MinInches`: Tamaño mínimo de pantalla (en pulgadas).
- `MaxInches`: Tamaño máximo de pantalla (en pulgadas).

### **Almacenamiento (SSD)**:
- `MinStorage`: Almacenamiento SSD mínimo (en GB).
- `MaxStorage`: Almacenamiento SSD máximo (en GB).

### **Precio**:
- `MinPrice`: Precio mínimo (en moneda local).
- `MaxPrice`: Precio máximo (en moneda local).

### **Procesador**:
- `Processor`: Modelo de procesador deseado (por ejemplo, `i5`, `Ryzen 7`).

---

## **Videos explicativos**

1. **Instalación y Ejecución**:
   [Ver en YouTube](https://www.youtube.com/watch?v=jMjJ0GVIAG4&ab_channel=franciscoSaez)

2. **Caso de Prueba**:
   [Ver en YouTube](https://www.youtube.com/watch?v=DANOlJklxLI&ab_channel=franciscoSaez)

---
## **Notas adicionales**

- Los datos obtenidos se procesan utilizando expresiones regulares (regex), lo que puede afectar la precisión en algunos casos.
- No siempre se logran extraer todos los atributos debido a la estructura de las páginas web objetivo.

---

Si tienes dudas o necesitas ayuda, por favor contacta a alguno de los integrantes.

