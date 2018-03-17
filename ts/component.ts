import * as Nightmare from 'nightmare'

export abstract class Component {
    url: string = ""
    partNum: string = ""
    brand: string = ""
    model: string = "" 
    rank: number = 0
    benchmark: number = 0
    samples: number = 0

    abstract isValid(old: Component): boolean
    abstract get(): Component
}

export class CPU extends Component {
    cores: string = ""
    averages: string[]
    performance: string[]
    subResults: string[]

    constructor() {
        super()
        this.averages = new Array(3)
        this.performance = new Array(3)
        this.subResults = new Array(9)
    }
    
    isValid(old: Component): boolean {
        return true
    }

    get(): CPU {

        return this
    }

    getCores(n: Nightmare) {
        
    }
}

export class GPU extends Component {
    constructor() {
        super()
        
    }

    isValid(old: Component): boolean {
        return true
    }

    get(): GPU {

        return this
    }
}