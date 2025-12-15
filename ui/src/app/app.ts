import { Component } from "@angular/core";
import { Packs } from "./packs/packs";
import { MatTab, MatTabGroup } from "@angular/material/tabs";
import { Package } from "./package/package";

@Component({
  standalone: true,
  selector: "app-root",
  imports: [
    Packs,
    MatTabGroup,
    MatTab,
    Package
  ],
  templateUrl: "./app.html",
  styleUrl: "./app.scss"
})
export class App {
}
