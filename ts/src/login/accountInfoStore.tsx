import { makeAutoObservable } from 'mobx';
import { createContext, useContext } from "react"

class AccountInfo {
    name: string = ''
    clan: string = ''
    permission: number = -1

    constructor() {
        makeAutoObservable(this);
    }
    setAccountName(name: string) {
        this.name = name;
    }

    getAccountName() {
        return this.name;
    }

    setClan(clan: string) {
        this.clan = clan;
    }
    getClan() {
        return this.clan;
    }

    setPermission(permission: number) {
        this.permission = permission;
    }
    getPermission() {
        return this.permission;
    }
}
const AccountInfoContext = createContext<AccountInfo>(new AccountInfo);

export { AccountInfoContext, AccountInfo };


