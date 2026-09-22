import {NewEC, NewT} from "./util.js";

export class DriverComponent {
    driverView = null;

    async render() {
        this.driverView = NewEC('div', 'driver-component');

        this.driverView.appendChild(NewT('TODO'));

        return this.driverView;
    }
}