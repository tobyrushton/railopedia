declare namespace journey {
    interface IStation {
        name: string
        code: string
    }

    interface IPrice {
        Provider: string
        Price: number
        Link: string
    }

    interface IJourney {
        ArrivalTime: string
        DepartureTime: string
        Prices: IJourneyPrice[]
    }

    interface IJourneyPrice extends Omit<IJourney, 'Prices'> {
        Prices: IPrice[]
    }
}
