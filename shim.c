#include <math.h>

void *func_acos  = acos;
void *func_asin  = asin;
void *func_atan  = atan;
void *func_cos   = cos;
void *func_cosh  = cosh;
void *func_sin   = sin;
void *func_sinh  = sinh;
void *func_tan   = tan;
void *func_tanh  = tanh;
void *func_exp   = exp;
void *func_log   = log;
void *func_log10 = log10;
void *func_sqrt  = sqrt;
void *func_fabs  = fabs;

double eval(void *code, double x, double y) {
	double (*func)(double, double) = code;
	return func(x, y);
}

void eval_2d(void *code, double *dst, double xmin, double xmax, int nx, double ymin, double ymax, int ny){
	int ix, iy;
	double x, y;
	double (*func)(double, double) = code;
	for(iy=0; iy<ny; iy++){
		y = ymin + ((ymax-ymin)*(iy+0.5))/ny;
		for(ix=0; ix<nx; ix++){
			x = xmin + ((xmax-xmin)*(ix+0.5))/nx;
			dst[iy*nx+ix] = func(x, y);
		}
	}
}

double call_func(void* f, double x){
	double (*func)(double) = f;
	return func(x);
}

