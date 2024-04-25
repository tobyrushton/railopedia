import csvToJson from 'convert-csv-to-json'
import fs from 'fs'

const stationListUrl = 'https://raw.githubusercontent.com/trainline-eu/stations/master/stations.csv'

const stationListCsv = await fetch(stationListUrl).then(res => res.text())

const rawJson = csvToJson.csvStringToJson(stationListCsv)

const stationListJSON = rawJson.map((station: any) => {
    if(station.country === 'GB')
        return {
            name: station.name,
            code: station.atoc_id,
        }
}).filter(station => station)

fs.writeFileSync('data/station-list.json', JSON.stringify(stationListJSON, null, 2))
