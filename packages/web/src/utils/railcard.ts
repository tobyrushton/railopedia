export type Railcard =
    | '16-25'
    | '26-30'
    | 'Senior'
    | 'Disabled'
    | 'F&F'
    | 'TT'
    | '16-17'
    | 'Veteran'
    | 'N'

export const railcards: Record<Railcard, string> = {
    '16-25': '16-25 Railcard',
    '26-30': '26-30 Railcard',
    Senior: 'Senior Railcard',
    Disabled: 'Disabled Railcard',
    'F&F': 'Family & Friends Railcard',
    TT: 'Two Together Railcard',
    '16-17': '16-17 Saver',
    Veteran: 'Veteran Railcard',
    N: 'No Railcard',
}
