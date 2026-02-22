import { Link } from 'react-router-dom'
import { ArrowRight, Star } from 'lucide-react'
import { Button } from '@/components/ui/button'

export function Hero() {
  return (
    <section className="relative min-h-[85vh] flex items-center overflow-hidden bg-gradient-to-br from-orange-50 via-white to-amber-50">
      {/* Decorative blobs */}
      <div className="absolute top-0 right-0 w-96 h-96 bg-orange-200/30 rounded-full -translate-y-1/2 translate-x-1/2 blur-3xl" />
      <div className="absolute bottom-0 left-0 w-72 h-72 bg-amber-200/20 rounded-full translate-y-1/2 -translate-x-1/2 blur-3xl" />

      <div className="container mx-auto max-w-6xl px-4 py-20 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Text */}
          <div className="space-y-6">
            {/* Social proof pill */}
            <div className="inline-flex items-center gap-2 rounded-full bg-orange-100 px-4 py-1.5 text-sm font-medium text-orange-700">
              <Star className="h-4 w-4 fill-orange-500 text-orange-500" />
              Rated 4.9 by 2,000+ happy customers
            </div>

            <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight text-gray-900 leading-tight">
              Cakes Made With{' '}
              <span className="text-primary relative">
                Love
                <svg
                  className="absolute -bottom-2 left-0 w-full"
                  viewBox="0 0 200 8"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M2 6C50 2 100 2 198 6"
                    stroke="currentColor"
                    strokeWidth="3"
                    strokeLinecap="round"
                  />
                </svg>
              </span>{' '}
              for Every Occasion
            </h1>

            <p className="text-lg text-muted-foreground max-w-md leading-relaxed">
              From birthday celebrations to dream weddings — order your perfect custom cake online
              and have it delivered fresh to your door.
            </p>

            <div className="flex flex-col sm:flex-row gap-3">
              <Button size="lg" asChild className="text-base">
                <Link to="/products">
                  Order Now
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild className="text-base">
                <Link to="/products">Browse Cakes</Link>
              </Button>
            </div>

            {/* Stats */}
            <div className="flex gap-8 pt-4">
              {[
                { label: 'Cakes Delivered', value: '15,000+' },
                { label: 'Happy Customers', value: '2,000+' },
                { label: 'Years of Baking', value: '12+' },
              ].map((stat) => (
                <div key={stat.label}>
                  <p className="text-2xl font-bold text-primary">{stat.value}</p>
                  <p className="text-xs text-muted-foreground">{stat.label}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Image */}
          <div className="relative flex justify-center">
            <div className="relative w-full max-w-md aspect-square">
              {/* Main image */}
              <div className="absolute inset-4 rounded-3xl overflow-hidden shadow-2xl">
                <img
                  src="https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=600"
                  alt="Beautiful decorated birthday cake"
                  className="h-full w-full object-cover"
                />
              </div>

              {/* Floating badge – top left */}
              <div className="absolute top-0 left-0 rounded-2xl bg-white shadow-lg p-3 flex items-center gap-2 text-sm font-semibold">
                <span className="text-2xl">🎂</span>
                <div>
                  <p className="text-gray-900">Custom Orders</p>
                  <p className="text-xs text-muted-foreground font-normal">Available daily</p>
                </div>
              </div>

              {/* Floating badge – bottom right */}
              <div className="absolute bottom-2 right-0 rounded-2xl bg-white shadow-lg p-3 text-sm">
                <div className="flex items-center gap-1">
                  {[1,2,3,4,5].map((i) => (
                    <Star key={i} className="h-3.5 w-3.5 fill-yellow-400 text-yellow-400" />
                  ))}
                </div>
                <p className="font-semibold text-gray-900 mt-0.5">4.9 / 5.0</p>
                <p className="text-xs text-muted-foreground">2k+ reviews</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
