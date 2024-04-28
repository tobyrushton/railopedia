import stationList from '../../../../data/station-list.json'

export const stationIsValid = (station: journey.IStation): boolean =>
    stationList.some((s) => s.code === station.code)
