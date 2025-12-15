import { Component, OnInit, signal } from "@angular/core";
import { FormArray, FormBuilder, FormControl, FormGroup, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatError, MatFormField, MatInput, MatLabel } from "@angular/material/input";
import { MatButton, MatIconButton } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { ApiService } from "../api.service";
import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpErrorResponse } from "@angular/common/http";
import { MatProgressSpinner } from "@angular/material/progress-spinner";
import { finalize } from "rxjs";

interface FormGroupValue {
  packs: FormArray<FormControl<number>>;
}

@Component({
  standalone: true,
  selector: "app-packs",
  imports: [
    MatFormField,
    MatInput,
    MatLabel,
    MatError,
    ReactiveFormsModule,
    MatIconButton,
    ReactiveFormsModule,
    MatIconModule,
    MatButton,
    MatIconButton,
    MatInput,
    MatProgressSpinner
  ],
  templateUrl: "./packs.html",
  styleUrl: "./packs.scss",
})
export class Packs implements OnInit {

  readonly loading = signal(false);
  readonly submitting = signal(false);

  readonly form: FormGroup<FormGroupValue>;

  constructor(
    private fb: FormBuilder,
    private readonly service: ApiService,
    private readonly toast: MatSnackBar
  ) {
    this.form = this.fb.group({
      packs: this.fb.array<FormControl<number>>([])
    });
  }

  ngOnInit(): void {
    this.loading.set(true);

    this.service.getPacks()
      .pipe(finalize(() => this.loading.set(false)))
      .subscribe({
        next: packs => {
          packs.forEach(pack => this.addPack(pack));
        },
        error: err => {
          this.toast.open(err?.error?.error ?? "Failed to sync packs", undefined, { duration: 5000 });
        }
      });
  }

  newPack(value: number | null = null): any {
    return this.fb.control(value, [Validators.required, Validators.min(1)]);
  }

  addPack(value: number | null = null) {
    this.form.controls.packs.push(this.newPack(value));
  }

  removeQuantity(i: number) {
    this.form.controls.packs.removeAt(i);
  }

  onSubmit() {
    this.submitting.set(true);

    const packs = this.form.value.packs ?? [];
    this.service.syncPacks(packs)
      .pipe(finalize(() => this.submitting.set(false)))
      .subscribe({
        next: () => {
        },
        error: (error: HttpErrorResponse) => {
          this.toast.open(error?.error?.error ?? "Failed to sync packs", undefined, { duration: 5000 });
        }
      });
  }
}
