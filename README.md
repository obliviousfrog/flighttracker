# Flight Tracker
A simple microservice that tracks a persons flight source and destination flights. 

## Installation

### Source
```
go build 
```

### Docker 
```
docker build -t flighttracker .
```

## Running
### Source
```
./flighttracker
```

### Docker
```
docker run -it -p 8080:8080 flightracker 
```

## Endpoints

### Calculate : `POST /calculate`


Used to calculate the begining and end flights of a flight list.

**URL** : `/calculate`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

Must be a double array/list with two values in at each index.

Good:
```json
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
```



Bad:

```json
[["IND", "EWR", "NYC"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
```

## Success Response

**Input**
```json
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
```

**Code** : `200 OK`

**Content example**

```json
["SFO", "EWR"]
```

## Error Response

**Condition** : If there is a problem parsing the data into a double list.

**Code** : `400 BAD REQUEST`

**Content**:
```json
"Error message."
```
<br>

**Condition** : If there is an issue sorting through the flight data.

**Code** : `500 BAD REQUEST`

**Content**:
```json
"Error message."
```
<br>

**Condition** : If there is an issue formating the result.

**Code** : `500 BAD REQUEST`

**Content**:
```json
"Error message."
```

