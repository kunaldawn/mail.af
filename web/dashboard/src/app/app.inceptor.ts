import {Injectable} from '@angular/core';
import {HttpEvent, HttpHandler, HttpInterceptor, HttpRequest,} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {Router} from "@angular/router";
import {catchError} from 'rxjs/operators';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
    constructor(private router: Router) {
    }

    intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        if (request.headers.has("X-Skip-Interceptor")) {
            const headers = request.headers.delete("X-Skip-Interceptor");
            return next.handle(request.clone({headers}));
        }

        let token = localStorage.getItem("af_token");
        console.log(token)
        if (token) {
            request = request.clone({
                setHeaders: {
                    Authorization: `Bearer ` + token
                }
            });
            return next.handle(request);
        }

        return next.handle(request).pipe(catchError(err => {
            if (err.status === 401) {
                this.router.navigate(["login"])
            }
            const error = err.error.message || err.statusText;
            return throwError(error);
        }));
    }
}
