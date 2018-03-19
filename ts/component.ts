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
    abstract get(body: HTMLBodyElement, n: Nightmare): Component
}

export class CPU extends Component {
    cores: string 
    averages: string[]
    subResults: string[]
    performance: string[]

    constructor() {
        super()
        this.cores = ""
        this.averages = new Array(3)
        this.performance = new Array(3)
        this.subResults = new Array(9)
    }
    
    isValid(old: CPU): boolean {
        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): CPU {
        this.getCores(body, n)
        this.getAverages(body, n)
        this.getSubResults(body)
        return this
    }

    getCores(body: HTMLBodyElement, n: Nightmare) {
        let sel = '.cmp-cpt.tallp.cmp-cpt-l'
        n.exists(sel)
            .then((exists: boolean) => {
                console.log("cores", exists)
                if(exists) {
                    this.cores = body.querySelector(sel)!.textContent!
                    console.log(this)
                    console.log(this.cores)
                    return
                }
                this.cores = ""
            })
    }

    getAverages(body: HTMLBodyElement, n: Nightmare) {
        for(var i = 0; i < 3; i++) {
            let selGreen  = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(${i+3}) .mcs-caption.pgbg`
            let selYellow = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(${i+3}) .mcs-caption.pybg`
            let selRed    = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(${i+3}) .mcs-caption.prbg`
            n.exists(`${selGreen}, ${selYellow}, ${selRed}`)
                .then((exists: boolean) => {
                    console.log("averages", exists)

                    if(exists){
                        this.averages[i] = body.querySelector(`${selGreen}, ${selYellow}, ${selRed}`)!.textContent!
                    }
                })
        }
    }

    getSubResults(body: HTMLBodyElement) {
        body.querySelectorAll('.mcs-hl-col').forEach((tag: Element, i: number, _: NodeListOf<Element>) => {
            this.subResults[i] = tag.textContent!
        }, this)
    }
}

export class GPU extends Component {
    constructor() {
        super()
        
    }

    isValid(old: GPU): boolean {
        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): GPU {

        return this
    }
}

export class RAM extends Component {
    constructor() {
        super()
    }

    isValid(old: RAM): boolean {

        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): RAM {
        
        return this
    }
}

export class SSD extends Component {
    constructor() {
        super()
    }

    isValid(old: SSD): boolean {

        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): SSD {

        return this
    }
}

export class HDD extends Component {
    constructor() {
        super()
    }

    isValid(old: HDD): boolean {
        
        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): HDD {

        return this
    }
}

export class USB extends Component {
    constructor() {
        super()
    }

    isValid(old: USB): boolean {

        return true
    }

    get(body: HTMLBodyElement, n: Nightmare): USB {

        return this
    }
}