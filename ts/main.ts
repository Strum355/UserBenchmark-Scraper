import * as Nightmare from 'nightmare'
import { CPU, GPU, RAM, SSD, HDD, USB, Component } from './component'

let nightmare = new Nightmare({show: true})

let c = new CPU()
c.url = "http://cpu.userbenchmark.com/Intel-Core-i7-8700K/Rating/3937"

try {
    nightmare
        .goto(c.url)
        .evaluate((): string => {
            return window.location.toString()
        }, (res: string) => {
            console.log(res)
        })
        .wait('body')
        .evaluate((): HTMLBodyElement => { 
            return document.querySelector('body')!
        }, (res: HTMLBodyElement) => {
            console.dir(res)
            console.log(res)
            c.get(res, nightmare)
        })
        .end()
        .then(() => {
            
        })
        .catch((e) => {
            console.log(e)
        })
} catch(e) {
    console.log(e)
}
