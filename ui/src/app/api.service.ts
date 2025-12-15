import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { environment } from "../environments/environment";
import { Observable } from "rxjs";

@Injectable({
  providedIn: "root",
})
export class ApiService {

  private readonly baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {
  }

  getPacks(): Observable<number[]> {
    return this.http.get<number[]>(`${this.baseUrl}/api/packs`);
  }

  syncPacks(values: number[]): Observable<any> {
    return this.http.post(`${this.baseUrl}/api/packs`, { packs: values });
  }

  package(value: number): Observable<Record<number, number>> {
    return this.http.post<Record<number, number>>(`${this.baseUrl}/api/package`, { numberOfItems: value });
  }
}
