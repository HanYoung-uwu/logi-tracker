import {makeAutoObservable} from 'mobx';
import {createContext, useContext} from "react"

class AccountInfo {
    name: string = ''

    constructor() {
        makeAutoObservable(this);
    }
    setAccountName(name:string) {
        this.name = name;
    }
}
const AccountInfoContext = createContext<AccountInfo>(new AccountInfo);

export {AccountInfoContext, AccountInfo};


