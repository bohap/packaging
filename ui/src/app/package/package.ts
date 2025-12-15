import { Component, signal } from "@angular/core";
import { FormControl, FormsModule, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatError, MatFormField, MatInput, MatLabel } from "@angular/material/input";
import { MatButton } from "@angular/material/button";
import { ApiService } from "../api.service";
import { HttpErrorResponse } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatProgressSpinner } from "@angular/material/progress-spinner";
import { finalize } from "rxjs";

@Component({
  selector: "app-package",
  imports: [
    MatFormField,
    FormsModule,
    MatError,
    MatInput,
    MatLabel,
    ReactiveFormsModule,
    MatButton,
    MatProgressSpinner
  ],
  templateUrl: "./package.html",
  styleUrl: "./package.scss",
})
export class Package {

  readonly control = new FormControl("", Validators.required);
  readonly result = signal<[number, number][] | null>(null);
  readonly submitting = signal(false);

  constructor(private readonly service: ApiService, private readonly toast: MatSnackBar) {
  }

  submit() {
    this.submitting.set(true);

    this.result.set(null);
    this.service.package(+(this.control.value!!))
      .pipe(finalize(() => this.submitting.set(false)))
      .subscribe({
        next: result => this.result.set(
          Object.entries(result).map(it => [+it[0], +it[1]])
        ),
        error: (error: HttpErrorResponse) => {
          this.toast.open(error?.error?.error ?? "Failed to package", undefined, { duration: 5000 });
        }
      });
  }
}
