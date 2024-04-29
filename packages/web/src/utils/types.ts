export const isIJourney = (
    arg: journey.IJourney[] | journey.IJourneyPrice[]
): arg is journey.IJourney[] => {
    if (
        (arg[0].Prices as journey.IPrice[]).every(
            (price) => price.Provider !== undefined && price.Price !== undefined
        ) &&
        arg.every(
            (journey) =>
                journey.ArrivalTime !== undefined &&
                journey.DepartureTime !== undefined
        )
    ) {
        return false
    } else return true
}
