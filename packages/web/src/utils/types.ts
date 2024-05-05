export const isIJourney = (
    arg: journey.IJourney[] | journey.IJourneyPrice[]
): arg is journey.IJourney[] => {
    return arg.length && !('Provider' in arg[0].Prices[0])
}
