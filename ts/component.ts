interface Component {
    readonly URL: string,
    readonly PartNum: string,
    readonly Brand: string,
    readonly Model: string,
    Rank: number,
    Benchmark: number,
    Samples: number,

    isValid(old: Component): boolean
}

class CPU implements Component {
    cores: string
    averages: string[]
    performance: string[]
    subresults: string[]
    URL: string
    PartNum: string
    Brand: string
    Model: string
    Rank: number
    Benchmark: number
    Samples: number

    constructor() {
        this.averages = new Array(3)
        this.performance = new Array(3)
        this.subresults = new Array(9)
    }
    
    isValid(old: Component): boolean {
        return true
    }
}