import { makeAutoObservable } from 'mobx';
import { createContext, useContext } from "react";
import { fetchLocations, Location } from '../api/apis';

class StockpileInfo {

    locations: Array<Location> = []

    constructor() {
        makeAutoObservable(this);
    }

    async refetchInfo() {
        fetchLocations().then((locations) => {
            if (locations) {
                this.setLocations(locations);
            }
        });
    }

    setLocations(locations: Array<Location>) {
        this.locations = locations;
    }
};

const LocationInfoContext = createContext<StockpileInfo>(new StockpileInfo);

export { LocationInfoContext, StockpileInfo };