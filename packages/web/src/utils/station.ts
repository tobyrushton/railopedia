import stationList from '../../../../data/station-list.json'

export const stationIsValid = (station: journey.IStation): boolean =>
    stationList.some((s) => s.code === station.code)

export const getStationName = (code: string): string => {
    const station = stationList.find((s) => s.code === code)
    return station ? station.name : ''
}
