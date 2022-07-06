import { makeAutoObservable } from 'mobx';
import { createContext, useContext } from "react"
import { ItemRecord, fetchAllItems } from "../api/apis";

class ItemsStore {
    
    items: Array<ItemRecord> = []

    constructor() {
        makeAutoObservable(this);
    }

    setItems(items: Array<ItemRecord>) {
        this.items = items;
    }

    getItems() {
        return this.items;
    }

    async refetchInfo() {
        let m_items = await fetchAllItems();
        if (m_items) {
            this.setItems(m_items);
        }
    }
}
const ItemsStoreContext = createContext<ItemsStore>(new ItemsStore);

export { ItemsStoreContext, ItemsStore };
