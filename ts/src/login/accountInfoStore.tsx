import { makeAutoObservable } from 'mobx';
import { createContext, useContext } from "react"

class AccountInfo {
    name: string = ''
    clan: string = ''
    permission: -2 | -1 | 0 | 1 | 2 | 3 | 4 = -1 // -2 unathorized, -1 undefined, 0 admin, 1 clan admin, 2 clan man,
    // 3 is temporary account for invitation links,
    // 4 is clan admin invitation links

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

    setPermission(permission: -2 | -1 | 0 | 1 | 2 | 3 | 4) {
        this.permission = permission;
    }
    getPermission() {
        return this.permission;
    }
}
const AccountInfoContext = createContext<AccountInfo>(new AccountInfo);

export { AccountInfoContext, AccountInfo };


